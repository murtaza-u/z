{
  description = "Murtaza Udaipurwala's Go Monolith";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-23.05";
  outputs = { self, nixpkgs, ... }@inputs:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      formatter.${system} = pkgs.nixpkgs-fmt;
      packages.${system} = {
        default = pkgs.buildGoModule {
          pname = "z";
          version = "0.2.0";
          src = ./.;
          vendorSha256 = "sha256-VBxDE+Ibalt8Po2oQbLS89hnC7M9VgKiUjV5v7PiySI=";
          CGO_ENABLED = 0;
          subPackages = [ "cmd/z" ];
          nativeBuildInputs = [ pkgs.installShellFiles ];
          postInstall = ''
            for shell in bash zsh; do
              installShellCompletion --$shell ./completion/$shell/z
            done
          '';
        };
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          go-tools
          gopls
        ];
      };
    };
}
