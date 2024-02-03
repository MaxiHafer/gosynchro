let stream = new EventSource("/gosynchro");

stream.addEventListener("gosynchro/connect", (event) => {
    console.log(`Connected to gosynchro stream on ${event.data["remote"]}`);
});

stream.addEventListener("gosynchro/reload", (event) => {
    console.log("received manual reload request");
    window.location.reload();
});

stream.addEventListener("gosynchro/filesystem", (event) => {
    console.log(`received filesystem event: ${JSON.stringify(event.data)}`);
    window.location.reload();
});