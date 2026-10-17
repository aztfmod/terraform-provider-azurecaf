import ast
import importlib.util
import os
from pathlib import Path
import subprocess
import tempfile
import unittest


SCRIPT_DIR = Path(__file__).resolve().parent
GENERATOR_PATH = SCRIPT_DIR / "generate_tests.py"

spec = importlib.util.spec_from_file_location("generate_tests", GENERATOR_PATH)
generate_tests = importlib.util.module_from_spec(spec)
assert spec.loader is not None
spec.loader.exec_module(generate_tests)


class GenerateTestsTest(unittest.TestCase):
    def test_dictionary_literals_do_not_contain_duplicate_keys(self):
        tree = ast.parse(GENERATOR_PATH.read_text())
        for node in ast.walk(tree):
            if not isinstance(node, ast.Dict):
                continue
            seen = {}
            for key in node.keys:
                if not isinstance(key, ast.Constant) or not isinstance(key.value, str):
                    continue
                self.assertNotIn(
                    key.value,
                    seen,
                    f"duplicate dictionary key {key.value!r} at lines "
                    f"{seen.get(key.value)} and {key.lineno}",
                )
                seen[key.value] = key.lineno

    def test_repeated_resource_overrides_are_merged(self):
        self.assertEqual(
            {
                "default_principals_modification_kind",
                "cluster_resource_id",
            },
            {
                key
                for key in generate_tests.RESOURCE_ATTR_OVERRIDES[
                    "azurerm_kusto_attached_database_configuration"
                ]
                if key in {"default_principals_modification_kind", "cluster_resource_id"}
            },
        )
        self.assertEqual(
            '"SafeNet Luna Network HSM A790"',
            generate_tests.RESOURCE_ATTR_OVERRIDES[
                "azurerm_dedicated_hardware_security_module"
            ]["sku_name"],
        )

    def test_resource_overrides_do_not_emit_mutually_exclusive_inputs(self):
        certificate_overrides = generate_tests.RESOURCE_ATTR_OVERRIDES[
            "azurerm_app_service_certificate"
        ]
        self.assertIn("key_vault_secret_id", certificate_overrides)
        self.assertNotIn("pfx_blob", certificate_overrides)
        self.assertNotIn("password", certificate_overrides)

        eventhub_overrides = generate_tests.RESOURCE_ATTR_OVERRIDES[
            "azurerm_iothub_endpoint_eventhub"
        ]
        self.assertIn("connection_string", eventhub_overrides)
        self.assertNotIn("endpoint_uri", eventhub_overrides)
        self.assertNotIn("entity_path", eventhub_overrides)

    def test_redhat_openshift_uses_a_valid_identity_fixture(self):
        identity_hcl = "\n".join(
            generate_tests.RESOURCE_EXTRA_HCL[
                "azurerm_redhat_openshift_cluster"
            ]
        )
        self.assertIn('type         = "UserAssigned"', identity_hcl)
        self.assertIn(
            "/providers/Microsoft.ManagedIdentity/userAssignedIdentities/uai-test",
            identity_hcl,
        )
        self.assertIn("platform_workload_identity_profile {", identity_hcl)
        self.assertIn("platform_workload_identity {", identity_hcl)

    def test_servicebus_disaster_recovery_uses_namespace_ids(self):
        overrides = generate_tests.RESOURCE_ATTR_OVERRIDES[
            "azurerm_servicebus_namespace_disaster_recovery_config"
        ]
        self.assertNotEqual(
            overrides["partner_namespace_id"],
            overrides["primary_namespace_id"],
        )
        for attribute in ("partner_namespace_id", "primary_namespace_id"):
            self.assertIn(
                "/providers/Microsoft.ServiceBus/namespaces/",
                overrides[attribute],
            )

    def test_sql_linked_services_use_passwordless_connection_strings(self):
        for resource_type in (
            "azurerm_data_factory_linked_service_azure_sql_database",
            "azurerm_data_factory_linked_service_sql_server",
        ):
            self.assertEqual(
                generate_tests.HCL_SQL_CONNECTION_STRING,
                generate_tests.RESOURCE_ATTR_OVERRIDES[resource_type][
                    "connection_string"
                ],
            )

    def test_logic_app_schema_override_is_valid_hcl(self):
        self.assertEqual(
            "jsonencode({})",
            generate_tests.RESOURCE_ATTR_OVERRIDES[
                "azurerm_logic_app_trigger_http_request"
            ]["schema"],
        )

    def test_resource_block_override_is_more_specific(self):
        block = {
            "block": {
                "attributes": {
                    "name": {"required": True, "type": "string"},
                    "tier": {"required": True, "type": "string"},
                }
            }
        }
        rendered = generate_tests.render_block(
            "sku", block, indent=2, resource_type="azurerm_application_gateway"
        )
        self.assertIn('name = "Standard_v2"', rendered)
        self.assertIn('tier = "Standard_v2"', rendered)


class RunAllTest(unittest.TestCase):
    def test_runner_isolates_init_and_requires_local_provider_override(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            workspace = root / "workspaces" / "azurerm_example"
            plugin_dir = root / "plugin"
            fake_bin = root / "bin"
            workspace.mkdir(parents=True)
            plugin_dir.mkdir()
            fake_bin.mkdir()

            provider = plugin_dir / "terraform-provider-azurecaf"
            provider.write_text("#!/usr/bin/env bash\nexit 0\n")
            provider.chmod(0o755)
            (workspace / "terraform.rc").write_text(
                generate_tests.make_terraform_rc(str(plugin_dir))
            )

            calls = root / "terraform-calls"
            terraform = fake_bin / "terraform"
            terraform.write_text(
                "#!/usr/bin/env bash\n"
                'printf "%s|%s\\n" "$1" "${TF_CLI_CONFIG_FILE:-}" >> "$CALL_LOG"\n'
                'if [[ "$1" == "test" ]]; then\n'
                '  echo "Success! 3 passed, 0 failed."\n'
                "fi\n"
            )
            terraform.chmod(0o755)

            report = root / "report.tsv"
            env = os.environ.copy()
            env.update(
                {
                    "CALL_LOG": str(calls),
                    "HOME": str(root),
                    "PATH": f"{fake_bin}{os.pathsep}{env['PATH']}",
                }
            )
            subprocess.run(
                [
                    str(SCRIPT_DIR / "run_all.sh"),
                    "--out-dir",
                    str(root / "workspaces"),
                    "--report",
                    str(report),
                ],
                check=True,
                capture_output=True,
                text=True,
                env=env,
            )

            init_call, test_call = calls.read_text().splitlines()
            self.assertEqual(
                f"init|{root / 'logs' / 'terraform-init.rc'}",
                init_call,
            )
            self.assertEqual(f"test|{workspace / 'terraform.rc'}", test_call)
            self.assertIn("\tPASS\t", report.read_text())
            self.assertIn(
                f"Using local azurecaf development override: {provider}",
                (root / "logs" / "azurerm_example.log").read_text(),
            )

    def test_runner_retries_transient_provider_startup_failures_until_success(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            workspace = root / "workspaces" / "azurerm_example"
            plugin_dir = root / "plugin"
            fake_bin = root / "bin"
            workspace.mkdir(parents=True)
            plugin_dir.mkdir()
            fake_bin.mkdir()

            provider = plugin_dir / "terraform-provider-azurecaf"
            provider.write_text("#!/usr/bin/env bash\nexit 0\n")
            provider.chmod(0o755)
            (workspace / "terraform.rc").write_text(
                generate_tests.make_terraform_rc(str(plugin_dir))
            )

            calls = root / "terraform-calls"
            attempt_count = root / "attempt-count"
            terraform = fake_bin / "terraform"
            terraform.write_text(
                "#!/usr/bin/env bash\n"
                'printf "%s|%s\\n" "$1" "${TF_CLI_CONFIG_FILE:-}" >> "$CALL_LOG"\n'
                'if [[ "$1" == "test" ]]; then\n'
                '  count=0\n'
                '  if [[ -f "$ATTEMPT_COUNT_FILE" ]]; then\n'
                '    count=$(cat "$ATTEMPT_COUNT_FILE")\n'
                '  fi\n'
                '  count=$((count + 1))\n'
                '  printf "%s\\n" "$count" > "$ATTEMPT_COUNT_FILE"\n'
                '  if [[ "$count" -lt 3 ]]; then\n'
                '    echo "Error: could not create provider instance: Unrecognized remote plugin message:"\n'
                '    echo "Failure! 0 passed, 0 failed."\n'
                "    exit 1\n"
                "  fi\n"
                '  echo "Success! 3 passed, 0 failed."\n'
                "fi\n"
            )
            terraform.chmod(0o755)

            report = root / "report.tsv"
            env = os.environ.copy()
            env.update(
                {
                    "CALL_LOG": str(calls),
                    "ATTEMPT_COUNT_FILE": str(attempt_count),
                    "HOME": str(root),
                    "PATH": f"{fake_bin}{os.pathsep}{env['PATH']}",
                }
            )
            subprocess.run(
                [
                    str(SCRIPT_DIR / "run_all.sh"),
                    "--out-dir",
                    str(root / "workspaces"),
                    "--report",
                    str(report),
                ],
                check=True,
                capture_output=True,
                text=True,
                env=env,
            )

            test_calls = [
                line
                for line in calls.read_text().splitlines()
                if line.startswith("test|")
            ]
            self.assertEqual(3, len(test_calls))
            self.assertIn("\tPASS\t3\t0\t", report.read_text())
            self.assertIn(
                "Transient provider startup failure detected; retrying attempt 3 of 3 after 1s.",
                (root / "logs" / "azurerm_example.log").read_text(),
            )


if __name__ == "__main__":
    unittest.main()
