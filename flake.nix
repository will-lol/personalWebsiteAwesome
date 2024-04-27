{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    templRepo.url = "github:a-h/templ/refs/tags/v0.2.543";
    aws-cdk = {
      url = "https://registry.npmjs.org/aws-cdk/-/aws-cdk-2.139.0.tgz";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, aws-cdk, templRepo }:
    flake-utils.lib.eachDefaultSystem (system: 
      let 
        name = "personalWebsite";
	templ = templRepo.packages.${system}.templ;
	cdk = pkgs.stdenv.mkDerivation {
	  name = "aws-cdk";
	  version = "2.121.1";
	  src = aws-cdk;
	  phases = [ "installPhase" "patchPhase" ];
	  installPhase = ''
	    cp -r $src $out
	  '';
	};
	overlays = [];
	lib = nixpkgs.lib;
	pkgs = import nixpkgs { inherit system overlays; };
	goBuild = (pkgs.buildGoModule {
	  inherit name;
	  src = ./.;
	  CGO_ENABLED = 0;
	  ldflags = ["-s -w"];
	  vendorHash = "sha256-54WH1NFO8fqVpINFWeBJZjh5Rn7H0lyOuiQclFP4Hx4=";
	  preBuild = ''
	    ${templ}/bin/templ generate
	  '';
	}).overrideAttrs (old: old // {GOOS = "linux"; GOARCH = "arm64"; });
      in
        {
	  packages = {
	    go = goBuild;
	  };
	  defaultPackage = goBuild;
	  devShell = pkgs.mkShell {
	    packages = [ pkgs.tailwindcss pkgs.entr pkgs.nodePackages.typescript-language-server pkgs.fd pkgs.nodejs_18 pkgs.awscli2 pkgs.aws-sam-cli pkgs.go pkgs.gopls templ cdk ];
	    shellHook = ''
	      export ENVIRONMENT=dev
	    '';
	  };
	}
    );
}
