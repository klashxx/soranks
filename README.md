## soranks

**S**tack**O**verflow **RAN**king**S** by *location*.

This piece of software is an exercise about consuming  [Stackexchange API](https://api.stackexchange.com).

Some time ago I wondered how to find SO *nuts* in my area.

My first natural approach was to use the query composer, find out the results [here](http://data.stackexchange.com/stackoverflow/query/381151/almeria-users).

Nice but very tied.

The second attempt was clear, direct consumption of  the `API`,  unfortunately [users method](https://api.stackexchange.com/docs/users) does not accept `location` filtering so I come with this code.

#### Examples

Usage:

```
$ ./soranks -h
Usage of ./soranks:
  -json string
    	json
  -limit int
    	max number of records (default 20)
  -location string
    	location (default "spain")
```

Any location:

```
$ ./soranks --location="."
Rank Name                              Rep Location
   1 Jon Skeet                      883667 Reading, United Kingdom
   2 Darin Dimitrov                 676217 Sofia, Bulgaria
   3 BalusC                         666857 Amsterdam, Netherlands
   4 Hans Passant                   640507 Madison, WI
   5 Marc Gravell                   615610 Forest of Dean, United Kingdom
   6 VonC                           606578 France
   7 CommonsWare                    575522 Up, Up, and Away
   8 SLaks                          526747 New Jersey
   9 Greg Hewgill                   496209 Christchurch, New Zealand
  10 Martijn Pieters                478812 Cambridge, United Kingdom
  11 Quentin                        474463 United Kingdom
  12 Alex Martelli                  464810 Sunnyvale, CA
  13 T.J. Crowder                   462818 United Kingdom
  14 JaredPar                       440782 Redmond, WA
  15 CMS                            440710 Guatemala
  16 marc_s                         440692 Bern, Switzerland
  17 Gordon Linoff                  440599 New York, United States
  18 dasblinkenlight                438190 United States
  19 Mark Byers                     435081 Denmark
  20 Guffa                          434711 Västervåla, Sweden
```

Specific location:

```
$ ./soranks --location="almer.a"
   1 José Juan Sánchez                7599 Almeria, Spain
   2 klashxx                          6034 Almeria, Spain
   3 Juanjo Vega                       799 Almeria, Spain
   4 segarci                           744 Huercal De Almeria, Spain
   5 Miguel Gil Martínez               411 Almeria, Spain
```
