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
            hash = "sha256-4mK86IJy/1FsG5ChW1jEdyz/QxSS0U0a+H6Kqd8CjZg=";
            vendorHash = "sha256-4mK86IJy/1FsG5ChW1jEdyz/QxSS0U0a+H6Kqd8CjZg=";
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
