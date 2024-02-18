# Rinha Backend v2

Vale ressaltar que esse código serve apenas para os objetivos da rinha... Portanto, nenhum padrão foi aplicado. Foi criado puramente por diversão. ·-·
Repo do zan -> [rinha backend 2ª edição](https://github.com/zanfranceschi/rinha-de-backend-2024-q1).

```
                                                                                                                       @@@@@@@@@@@@@@@ @@@@@@@@         
                                                                         @@@@@          @@@@@                         @@@@@@@@@@@@@@@@@@@@@@@@@@        
                                                                         @@@@@@         @@@@@                        @@@@@@@@@@@@@@@@@@@@@@@@@@@@       
                                                                         @@@@@@@@       @@@@@                        @@@@@@@@@@@@@@@@@@@@@  @@@@@       
                                                                         @@@@@@@@@@     @@@@@                        @@@@@@@@@@@@@@@@@@@@@@ @@@@        
                      @@@@@@@@@@@     @@@@@@@@@@@                        @@@@ @@@@@@@   @@@@@                         @@@@@@@@@@@@@@@@@@@@@@@@@@        
                    @@@@@@@@@@@@@@  @@@@@@@@@@@@@@@                      @@@@   @@@@@@@ @@@@@                         @@@@@@@@@@@@@@@@@@@ @ @@@         
                   @@@@@@     @@@  @@@@@@     @@@@@@                     @@@@     @@@@@@@@@@@                          @@@@@@@@@@@@@@@@@@@ @@@          
      @@@@@@@@@@  @@@@@   @@@@@@@@@@@@@        @@@@@                     @@@@       @@@@@@@@@                           @@@@@@@@@  @@@@@@@   @ @@       
           @@@@@  @@@@@   @@@@@@@@@@@@@        @@@@@                     @@@@         @@@@@@@                            @@@@@@@@@ @@@@@@@@@@@@@        
                  @@@@@      @@@@@@@@@@       @@@@@                      @@@@           @@@@@                             @@@@@@@@@@@@@@@@@             
                  @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                        @@              @                                       @@@@@@@ @             
                    @@@@@@@@@@@@    @@@@@@@@@@@@@              @@@    @@  @@@@@@  @@ @@    @@ @@    @@                            @@@@@@@@@             
                       @@@@@@          @@@@@@                  @@@@@  @@ @@    @@ @@ @@  @@@@  @@@@@@                              @@@@ @@              
                                                               @@@ @@@@@ @@ @@@@@ @@@@@@@@ @@   @@@@              @@@        @       @   @@@ @@@@ @     
                                                               @@@   @@@  @@@@@@  @@@@@@   @@ @@@  @@@            @@@@@@@@@@ @@@@@@@@@@@@ @@@@  @ @     
                                                                                                                  @   @ @@ @@@ @@@@@ @@    @ @  @ @     
                                                                                                                               @@@                      
```

## Stack

* Golang (v1.22)
* PostgreSQL
* NGINX

## Usage

```zsh
$ make

Usage: make <target>
  up                         Build containers
  clean                      Clear all
  health.check               Check if it went up ok
  stress                     Run stress tests
  docker.build               Build docker image
  docker.push                Push docker image
```

## Run

```zsh
$ docker-compose -f docker-compose.yml up -d --build
# Or
$ make up
```

Health-check
```zsh
$ curl -v http://localhost:9999
```

## Stress-test
Obs: You can see the preview in result (run on my machine) ·-·
```zsh
$ make stress 
# Maybe this command doesn't work for you
$ google-chrome-stable stress-test/user-files/results/**/index.html
```

#### That's it, so... Shinzo wo sasageyo