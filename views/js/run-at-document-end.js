import { render, html } from './uhtml/index.js'

for (const el of document.querySelectorAll('.artwork-container,.artwork-container-scroll')) {
  const artwork_container_actions_node = document.createElement('div');
  el.parentElement.insertBefore(artwork_container_actions_node, el);
  const openAllArt = () => {
    let xs = el.querySelectorAll('div.artwork-small > a')
    Array.from(xs).map(x=> window.open(x.href))
  }
  render(artwork_container_actions_node, html`<button onclick=${openAllArt}>Open all artworks</button>`);
}
