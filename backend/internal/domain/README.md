# Services

## Main

### Playlist
    1. Meta
        depends on:
            - Policy
            - PlaylistMetaRepository

    2. Tracks
        depends on:
            - Policy
            - PlaylistTracksRepository

    3. Favorites
        depends on:
            - Policy
            - PlaylistFavoritesRepository

    4. Deliter
        depends on:
            - Policy -?
            - MetaDeleter
            - TracksDeleter
            - FavoritesDeleter

    5. Policy
        depends on:
            - PlaylistAccessRepository 
###