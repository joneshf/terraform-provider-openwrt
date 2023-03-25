{
  description = "A Terraform provider for OpenWrt";

  inputs = {
    nixpkgs = {
      owner = "NixOS";
      ref = "nixos-22.11";
      repo = "nixpkgs";
      type = "github";
    };
  };

  outputs = { self, nixpkgs }:
    let
      # Helper to provide system-specific attributes
      forAllSupportedSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });

      supportedSystems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
    in

    {
      devShells = forAllSupportedSystems ({ pkgs }: {
        default = pkgs.mkShell {
          packages = [
            pkgs.colima
            pkgs.docker
            pkgs.docker-credential-helpers
            pkgs.gh
            pkgs.git
            pkgs.gnumake
            pkgs.gnupg
            pkgs.go_1_20
            pkgs.go-tools
            pkgs.gopls
            pkgs.goreleaser
            pkgs.jq
            pkgs.nixpkgs-fmt
            pkgs.semver
            pkgs.terraform
          ];
        };
      });
    };
}
