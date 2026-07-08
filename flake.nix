{
  description = "Tailscale SSH ProxyCommand helper";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      systems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.buildGoModule {
            pname = "ts-proxy-command";
            version = "0.1.0";

            src = ./.;
            vendorHash = "sha256-5zGexM5aiw09cCgzLzX3ulRXsK5w3efIpL6KeleHQak=";

            meta = {
              description = "Tailscale tsnet SSH ProxyCommand helper";
              mainProgram = "ts-proxy-command";
            };
          };
        }
      );

      checks = forAllSystems (system: {
        default = self.packages.${system}.default;
      });
    };
}
