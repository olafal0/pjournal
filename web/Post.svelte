<script>
  import marked from "marked";
  import hljs from "highlight.js";

  marked.setOptions({
    gfm: true,
    highlight: function(code, language) {
      const validLanguage = hljs.getLanguage(language) ? language : "plaintext";
      return hljs.highlight(validLanguage, code).value;
    }
  });

  export let content = "";
  export let createdAt = Date.now();
  let timestamp;

  $: timestamp = new Date(createdAt).toLocaleString(undefined, {
    weekday: "short",
    month: "long",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "numeric"
  });
</script>

<style>
  .aside {
    font-size: 75%;
  }
</style>

<div class="row">
  <div class="col l8 s12 offset-l2">
    <div class="card grey darken-4">
      <div class="card-content white-text">
        <i class="aside">{timestamp}</i>
        <p>
          {@html marked(content, { gfm: true })}
        </p>
      </div>
    </div>
  </div>
</div>
