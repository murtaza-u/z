{
  description = "Go Monolith";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.05";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { self, nixpkgs, ... }@inputs:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ inputs.gomod2nix.overlays.default ];
      };
    in
    {
      formatter.${system} = pkgs.nixpkgs-fmt;
      packages.${system} = {
        default = pkgs.buildGoApplication {
          pname = "z";
          version = "0.2.0";
          src = ./.;
          modules = ./gomod2nix.toml;
          subPackages = [ "cmd/z" ];
        };
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          go-tools
          gopls
          inputs.gomod2nix.packages.${system}.default
        ];
      };
    };
}
