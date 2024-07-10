window.onload = () => {
    document.body.style.overflow = 'hidden';
    document.querySelectorAll(".playlist-container").forEach(element => element.addEventListener("click", () => {
        const choosen_playlist_id = element.children[2].innerHTML;
        
        const new_url = `/popularity?playlist=${choosen_playlist_id}`
        window.location.href = new_url;

    }))

    
};



window.onSpotifyWebPlaybackSDKReady = async () => {
    const token = await (await fetch("/token")).text();
    const player = new Spotify.Player({
      name: 'Rate My Playlist WebApp',
      getOAuthToken: cb => { cb(token); },
      volume: 0.5
    });
      // Ready
    player.addListener('ready', ({ device_id }) => {
        console.log('Ready with Device ID', device_id);
    });

    // Not Ready
    player.addListener('not_ready', ({ device_id }) => {
        console.log('Device ID has gone offline', device_id);
    });


    player.addListener('initialization_error', ({ message }) => {
        console.error(message);
    });
  
    player.addListener('authentication_error', ({ message }) => {
        console.error(message);
    });
  
    player.addListener('account_error', ({ message }) => {
        console.error(message);
    });


    player.connect();

    const song_toggle_button = document.getElementById("music-control-button");
    console.log(song_toggle_button)


    song_toggle_button.addEventListener("click", (e) =>{
        console.log("asd")
        e.preventDefault();
        player.togglePlay();
    })
}

