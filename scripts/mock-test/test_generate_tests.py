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


if __name__ == "__main__":
    unittest.main()
