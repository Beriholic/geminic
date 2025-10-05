{
  description = "Geminic";

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
        pkgs = import nixpkgs { inherit system; };

        pname = "geminic";
        version = "0.5.0";
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            inherit pname version;
            src = ./.;
            hash = "sha256-v2Xcfm582FBiG1ZZCSAF6MilJ0YOTf7ozLv81LO/Xjk=";
            vendorHash = "sha256-v2Xcfm582FBiG1ZZCSAF6MilJ0YOTf7ozLv81LO/Xjk=";
          };
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.go
              pkgs.gopls
              pkgs.delve
              pkgs.gotools
            ];
          };
        };

        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/${pname}";
          };
        };
      }
    );
}
