<template>
  <div>
    <input v-model="q" @input="load" placeholder="Search states..." />
    <ul v-if="list.length">
      <li v-for="s in list" :key="s.code" @click="choose(s)">{{ s.name }}</li>
    </ul>
  </div>
</template>
<script>
import { ref } from "vue";
export default {
  emits: ["select"],
  setup(_, { emit }) {
    const q = ref("");
    const list = ref([]);
    async function load() {
      if (!q.value) {
        list.value = [];
        return;
      }
      const res = await fetch(import.meta.env.VITE_API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: "query($q:String!){states(q:$q){name code}}",
          variables: { q: q.value },
        }),
      });
      list.value = (await res.json()).data.states;
    }
    function choose(s) {
      emit("select", s);
      q.value = s.name;
      list.value = [];
    }
    return { q, list, load, choose };
  },
};
</script>
<style>
ul {
  list-style: none;
  padding: 0;
  border: 1px solid #ccc;
}
li {
  padding: 4px;
  cursor: pointer;
}
li:hover {
  background: #eef;
}
</style>
