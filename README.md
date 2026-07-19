# csv2ynab

Parses CSV and Excel files from some Belgian banks and either dumps them into a CSV that YNAB can read or post to the YNAB API directly.

## Usage

```bash
go install github.com/victorjacobs/csv2ynab@latest
csv2ynab -in export.csv
```

The flake also exposes a package, an app, and a NixOS module. A minimal NixOS
configuration looks like this:

```nix
{
  inputs.csv2ynab.url = "github:victorjacobs/csv2ynab";

  outputs = { nixpkgs, csv2ynab, ... }: {
    nixosConfigurations.example = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        csv2ynab.nixosModules.default
        {
          services.csv2ynab = {
            enable = true;
            bindAddress = "127.0.0.1:8080";
            ynabAPIKeyFile = "/run/secrets/ynab-api-key";
            settings.watch.patterns = [
              {
                pattern = "*.csv";
                budget_id = "your-budget-id";
                account_id = "your-account-id";
              }
            ];
          };
        }
      ];
    };
  };
}
```

The API key file is loaded through systemd credentials and is only added to the
generated configuration under `/run`; it is not copied into the Nix store.
