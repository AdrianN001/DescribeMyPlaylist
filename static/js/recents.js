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

        return {
            artists: artists,
            title: track.name,
            preview_url: track.preview_url
        }
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

    // let audio = new Audio(json_data["preview_url"])
    // audio.type = "audio/mp3";
    // audio.play().then(x => console.log(x)).catch(e => console.error(e))

    const [artists, tracks] = format_incomming_json_data(json_data);

    const artist_container = document.querySelector("#artists");

    while(artist_container.firstChild){
        artist_container.removeChild(artist_container.lastChild);
    }

    artists.forEach(artist => {
        const name = artist.name;
        const image = artist.image;


        const new_tag = document.createElement("div")
        new_tag.className = "artist"
        new_tag.innerHTML = `<img class="artist_image" src="${image}"/> <div class="artist_name">${name}</div>`
        artist_container.appendChild(new_tag)


    })

    const track_container = document.querySelector("#tracks");

    while(track_container.firstChild){
        track_container.removeChild(track_container.lastChild);
    }
    
    tracks.forEach(track => {
        const title = track.title;
        const artists = track.artists;
        const preview_url = track.preview_url;


        const new_tag = document.createElement("div")
        new_tag.className = "track"
        new_tag.innerHTML = `${artists} - <span class="track_name">${title}</span>`
        track_container.appendChild(new_tag)


        new_tag.addEventListener('click', (e) => {
            e.preventDefault()

            // audio.pause();
            // audio = new Audio(preview_url)
            // audio.type = "audio/mp3";
            // audio.play().then(x => console.log(x)).catch(e => console.error(e))

        })
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

function debounce(fn, delay) {
    let timer;
    return (() => {
      clearTimeout(timer);
      timer = setTimeout(() => fn(), delay);
    })();
    
};


document.getElementById("long_term_button").addEventListener('click', (e) => update_ui_by_new_timerange(TIME_RANGES.LONG_TERM))
document.getElementById("mid_term_button").addEventListener('click', (e) => update_ui_by_new_timerange(TIME_RANGES.MID_TERM))
document.getElementById("short_term_button").addEventListener('click', (e) => update_ui_by_new_timerange(TIME_RANGES.SHORT_TERM))
