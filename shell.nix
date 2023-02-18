let

  nixpkgs-tarball_22-11 = builtins.fetchTarball {
    sha256 = "sha256:11w3wn2yjhaa5pv20gbfbirvjq6i3m7pqrq2msf0g7cv44vijwgw";
    # 22.11
    url = "https://github.com/NixOS/nixpkgs/archive/4d2b37a84fad1091b9de401eb450aae66f1a741e.tar.gz";
  };

in

{ pkgs ? import nixpkgs-tarball_22-11 { } }:

pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.gh
    pkgs.git
    pkgs.go
    pkgs.go-tools
    pkgs.nixpkgs-fmt
    pkgs.terraform
  ];
}
