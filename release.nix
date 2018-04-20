{ pkgs ? import <nixpkgs> {} }:
with pkgs;
rec {
  repost = callPackage ./default.nix {};
  docker = dockerTools.buildImage {
    name = "jahkeup/repost";
    contents = repost;
    config = {
      Cmd = [ "${repost}/bin/repostd" ];
    };
  };
}