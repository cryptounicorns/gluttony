with import <nixpkgs>{};
{ pkgs ? import <nixpkgs> {} }:

buildGo19Package rec {
  name = "gluttony-unstable-${version}";
  version = "development";

  buildInputs = with pkgs; [ git glide ];

  src = ./.;
  goPackagePath = "github.com/cryptounicorns/gluttony";
}
