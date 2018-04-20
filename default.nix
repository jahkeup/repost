{ buildGoPackage, go }:
buildGoPackage {
  name = "repost";
  version = "0.1.0";

  src = ./.;

  goPackagePath = "github.com/jahkeup/repost";
}
