# cross-watchlist
cross-watchlist was initially just a watchlist matcher, but i'm slowly adding new funcs based on my needs:
random from trending, random from lists, random from watchlists, random based in profile. 

<!-- ![Output](https://i.imgur.com/3jAtq3M.png) -->

TODO:

	- Apply random to list and return a favorite one
	- showRandomFromList() (ok), maybe add cache
	- showRandomFromTrending()
	- showSuggestionsBasedOnProfile()
	- showCrossWatchlist() //Maybe the user can add preferences, length, genre, etc

    - update frontend for new funcs
    - update FE for renaming (submit-form to /random-from-watchlists)

QOL: 

    - Update frontend for mobile
    - Make "logo" in the side bar
    - move sidebar to top? 

    - add functions to the new modes
    - images in result
    - search a vps (no need for db)
    - It literally doesn't check for errors
        - IT WILL BREAK if you try to use: empty watchlists, invalid users!
    