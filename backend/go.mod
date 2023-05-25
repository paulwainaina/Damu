module example.com/damu

go 1.20

replace example.com/patient => ./modules/patient

replace example.com/triage => ./modules/triage

require (
	example.com/patient v0.0.0-00010101000000-000000000000
	example.com/triage v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
)
