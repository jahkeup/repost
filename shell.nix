{ pkgs ? import <nixpkgs> {} }:
with pkgs;
let
  dep = buildGoPackage {
    name = "golang-dep";
    src = fetchFromGitHub {
      owner = "golang";
      repo = "dep";
      rev = "d5c4d780bdd70faf9bf8574f704976bce465aaf1";
      sha256 = "1k8fa66252fgz44yd5xi8x9zldw6bpv9j3l4bjrph0mvm9gsy6gg";
    };
    goPackagePath = "github.com/golang/dep";
    subPackages = [ "cmd/dep" ];
  };
in
stdenv.mkDerivation rec {
  name = "repost-shell";
  buildInputs = [ go dep ];
}
