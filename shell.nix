{ pkgs ? import <nixpkgs> {} }:
with pkgs;
let
  gomock = buildGoPackage {
    name = "mockgen";
    # src = fetchFromGitHub {
    #   owner = "golang";
    #   repo = "mock";
    #   rev = "8b2eeeb0ca5f56c78bec5efde9c4a21d9201126c";
    #   sha256 = "1ldrxvmdr6sbhf88jvqzjw33222agaf5bnla9d70v1n4rsvfhyh8";
    # };
    src = ./vendor/github.com/golang/mock;
    goPackagePath = "github.com/golang/mock";
    subPackages = [ "mockgen" ];
    postConfigure = ''
    # Replace the older context references, its provided by the language now.
    find go/src/$goPackagePath -name \*.go \
         -exec sed -i 's,golang.org/x/net/context,context,g' {} \;
    '';
  };
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
  buildInputs = [ go dep gomock ];
}
