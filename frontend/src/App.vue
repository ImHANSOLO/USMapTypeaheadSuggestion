<template>
  <div id="app">
    <Typeahead @select="highlightState" />
    <div id="map" style="height: 500px; margin-top: 1em"></div>
  </div>
</template>
<script>
import Typeahead from "./components/Typeahead.vue";
export default {
  components: { Typeahead },
  data() {
    return { map: null };
  },
  methods: {
    highlightState(s) {
      this.map.data.setStyle((f) => ({
        fillColor: f.getProperty("STATE") === s.code ? "#ff0000" : "#ffffff",
        strokeWeight: 1,
      }));
    },
  },
  mounted() {
    window.initMap = () => {
      this.map = new google.maps.Map(document.getElementById("map"), {
        center: { lat: 37.8, lng: -96 },
        zoom: 4,
      });
      this.map.data.loadGeoJson(
        "https://raw.githubusercontent.com/PublicaMundi/MappingAPI/master/data/geojson/us-states.json"
      );
    };
    const tag = document.createElement("script");
    tag.src = `https://maps.googleapis.com/maps/api/js?key=${
      import.meta.env.VITE_GOOGLE_MAPS_API_KEY
    }&callback=initMap`;
    tag.async = true;
    document.head.appendChild(tag);
  },
};
</script>
<style>
ul {
  list-style: none;
  padding: 0;
  border: 1px solid #ccc;
  max-width: 300px;
}
li {
  padding: 4px;
  cursor: pointer;
}
li:hover {
  background: #eef;
}
</style>
