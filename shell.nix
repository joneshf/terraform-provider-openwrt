let

  nixpkgs-tarball_22-11 = builtins.fetchTarball {
    sha256 = "sha256:0745rigamnnzz4qf712pvjs3vz8qsg3r9g903k6m4z92yxr1w942";
    # 22.11
    url = "https://github.com/NixOS/nixpkgs/archive/e6d5772f3515b8518d50122471381feae7cbae36.tar.gz";
  };

in

{ pkgs ? import nixpkgs-tarball_22-11 { } }:

pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.colima
    pkgs.docker
    pkgs.docker-credential-helpers
    pkgs.gh
    pkgs.git
    pkgs.gnumake
    pkgs.go
    pkgs.go-tools
    pkgs.gopls
    pkgs.nixpkgs-fmt
    pkgs.terraform
  ];
}
