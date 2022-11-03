export default {
    props: {
      repo: Object
    },
    template: `
    <div class="card">
    <div class="card-body">
      <h5 class="card-title"><a :href="repo.url" target="_blank">{{ repo.name }}</a></h5>
      <p class="card-text">{{ repo.description }}</p>
      <a class="btn btn-danger">Abandon</a>
    </div>
  </div>
    `
  }