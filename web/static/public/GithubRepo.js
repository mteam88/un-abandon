export default {
    props: {
      repo: Object,
      buttontext: String,
    },
    template: `
    <div class="card" style="column-break-inside: avoid;">
    <div class="card-body">
      <h5 class="card-title"><a :href="repo.html_url" target="_blank">{{ repo.name }}</a></h5>
      <p class="card-text">{{ repo.description }}</p>
      <a class="btn btn-danger" @click="$emit('button-clicked')" >{{ buttontext }}</a>
    </div>
  </div>
    `
  }