{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = rec {
          csv2ynab = pkgs.callPackage ./nix/package.nix { };
          default = csv2ynab;
        };

        apps.default =
          flake-utils.lib.mkApp {
            drv = self.packages.${system}.default;
          }
          // {
            meta.description = "Import Belgian bank exports into YNAB";
          };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_27
            gopls
            gotools
            nixfmt
          ];
        };
      }
    )
    // {
      nixosModules = rec {
        csv2ynab = import ./nix/module.nix;
        default = csv2ynab;
      };
    };
}
