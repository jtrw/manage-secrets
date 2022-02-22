module github.com/boltdb/bolt

go 1.17

require (
	github.com/nilBora/bolt v1.3.1 // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	store/bolt v1.0.0 // indirect
)


replace store/bolt v1.0.0 => ./cmd/backend/app/store