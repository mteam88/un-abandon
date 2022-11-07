export default {
    methods: {
        addAlert: function (type, message) {
            var alert = document.createElement('div');
            alert.innerHTML =
                `<div class="alert alert-${type} alert-dismissible fade show" role="alert">
                    ${message}
                    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
                </div>`;
            this.$refs.alerts.appendChild(alert);
        },
    },
    template: `
    <div ref="alerts"></div>
    `
  }