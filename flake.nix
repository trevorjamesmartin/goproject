{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url  = "github:cachix/devenv";
    templ.url   = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, devenv, systems, ... } @ inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
      templ = system: inputs.templ.packages.${system}.templ;
    in
    {
      packages = forEachSystem (system: {
        devenv-up = self.devShells.${system}.default.config.procfileScript;

        goproject = nixpkgs.legacyPackages.${system}.buildGoModule {
          preBuild = ''
            if [[ ! -f ./static/htmx.min.js ]]; then
              echo "downloading htmx.min.js"
              curl -o $HTMXMIN ./static/htmx.min.js https://unpkg.com/htmx.org@1.9.11/dist/htmx.min.js
            fi
            ${templ system}/bin/templ generate
          '';
        };

      });

      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                {

                  env.PGUSER = "postgres";
                  env.PGPASSWORD = "secret";
                  
                  env.DBNAME = "appdata";
                  env.PWFILE = "/tmp/password.txt";

                  languages.javascript.enable = true;
                  languages.javascript.package = pkgs.nodejs;
                  
                  services.postgres.listen_addresses = "127.0.0.1";
                  services.postgres.port = 5432;
                  services.postgres.enable = true;

                  languages.go.enable = true;
                  difftastic.enable = true;

                  # https://devenv.sh/reference/options/
                  packages = with pkgs; [
                    postgresql
                    go
                    (templ system)
                    tree-sitter
                  ];

                  enterShell = ''
                    # exit gracefully 
                    trap gracefully 1 2 3 6 15 EXIT

                    gracefully ()
                    {
                      echo "gracefully exiting..."
                      pg_ctl stop
                    }

                    has-htmx()
                    {
                      if [[ ! -f ./static/htmx.min.js ]]; then
                        echo "downloading htmx.min.js"
                        curl -o ./static/htmx.min.js https://unpkg.com/htmx.org@1.9.11/dist/htmx.min.js
                        curl -o ./static/class-tools.js https://unpkg.com/htmx.org@1.9.11/dist/ext/class-tools.js
                        curl -o ./static/reaniebeaniev20.woff2 https://fonts.gstatic.com/s/reeniebeanie/v20/z7NSdR76eDkaJKZJFkkjuvWxXPq1q6Gjb_0.woff2
                      fi
                    }

                    # ensure the environment exists
                    environment-exists ()
                    {
                      if [[ ! -d $PGDATA/logs ]]; then
                        echo "iniitializing database ..."
                        
                        echo $PGPASSWORD > $PWFILE                        

                        initdb \
                          --username=$PGUSER \
                          --encoding=UTF-8 \
                          --locale=en_US.UTF-8 \
                          --auth=md5 \
                          --pwfile=$PWFILE

                        mkdir -p $PGDATA/{run,logs}
                        
                        echo "updating postgresql.conf"
                        sed -i -e "s@#unix_socket_directories *= *.*@unix_socket_directories = '$PGDATA/run'@" $PGDATA/postgresql.conf
                        echo "...initialized"
                        go get github.com/a-h/templ
                        go get github.com/lib/pq
                        go mod tidy
                      fi
                    }

                    # start the database server
                    start-db ()
                    {
                      if [[ $(pg_ctl status|awk '{ print $2 }') == "no" ]]; then
                        pg_ctl -D $PGDATA -l $PGDATA/logs/postgres.log start
                      fi

                      # create application database
                      if [[ -f $PWFILE ]]; then
                        echo creating database: $DBNAME
                        createdb -U $PGUSER $DBNAME
                        rm $PWFILE
                      fi
                    }

                    environment-exists
                    
                    has-htmx
                    
                    if [[ ! -f .envrc ]]; then
                      # start the db server now
                      start-db
                    fi

                  '';

                }
              ];
            };
          });
    };
}


