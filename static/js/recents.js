const TIME_RANGES = Object.freeze({
    LONG_TERM: "long_term",
    MID_TERM:  "medium_term",
    SHORT_TERM:"short_term"
})

function format_incomming_json_data(json_data) {
    const tracks = json_data["tracks"].map(track => {
        
        const artists = track.artists.reduce((acc, artist) => {
            if (acc === ""){
                return artist.name
            }
            return `${acc}, ${artist.name}`
        }, "")

        return `${artists} - ${track.name}`
    });

    const artists = json_data["artists"].map(artist => {
        return {
            name: artist.name,
            image : artist.images[0].url
        }
    })

    return [artists, tracks]
}

function rendering_callback(json_data){
    console.log(json_data)

    const audio = new Audio(json_data["preview_url"])
    audio.type = "audio/mp3";
    audio.play().then(x => console.log(x)).catch(e => console.error(e))

    const [artists, tracks] = format_incomming_json_data(json_data);

    const artist_container = document.querySelector("#artists");
    artists.forEach(artist => {
        const name = artist.name;
        const image = artist.image;


        const new_tag = document.createElement("div")
        new_tag.className = "artist"
        new_tag.innerHTML = `<img class="artist_image" src="${image}"/> <div class="artist_name">${name}</div>`
        artist_container.appendChild(new_tag)

    })

    const track_container = document.querySelector("#tracks");
    tracks.forEach(track => {
        const name = track;


        const new_tag = document.createElement("div")
        new_tag.className = "track"
        new_tag.innerHTML = name
        track_container.appendChild(new_tag)

    })
}



function update_ui_by_new_timerange(time_range){
    fetch(`/get_recents?time_range=${time_range}`)
                                .then(val => val.json())
                                .then(rendering_callback)
                                .catch(err => console.error(err))
}


/* Initial run */ 
update_ui_by_new_timerange(TIME_RANGES.LONG_TERM)