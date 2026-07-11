<script>
  let { open = false } = $props();
</script>

<!-- Eye of Sauron: geometric low-poly style. open=connected, closed=reconnecting -->
<svg class="sauron-eye" class:open
     xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40" width="28" height="28"
     aria-label={open ? 'Connected' : 'Reconnecting…'} role="img">
  <!-- Eye body — scaleY squashes to simulate open/blink/closed -->
  <g class="eye-body">

    <!-- 24 triangular shards, alternating deep-red / orange -->
    <polygon points="18.5,3.1 21.5,3.1 20.0,15.0"    fill="#bb0000"/>
    <polygon points="22.9,3.2 25.8,4.0 21.3,15.2"    fill="#ff4400"/>
    <polygon points="27.2,4.6 29.7,6.1 22.5,15.7"    fill="#bb0000"/>
    <polygon points="30.9,7.0 33.0,9.1 23.5,16.5"    fill="#ff4400"/>
    <polygon points="33.9,10.3 35.4,12.8 24.3,17.5"  fill="#bb0000"/>
    <polygon points="36.0,14.2 36.7,17.0 24.8,18.7"  fill="#ff4400"/>
    <polygon points="36.9,18.5 36.9,21.5 25.0,20.0"  fill="#bb0000"/>
    <polygon points="36.7,23.0 36.0,25.8 24.8,21.3"  fill="#ff4400"/>
    <polygon points="35.4,27.2 33.9,29.7 24.3,22.5"  fill="#bb0000"/>
    <polygon points="33.0,30.9 30.9,33.0 23.5,23.5"  fill="#ff4400"/>
    <polygon points="29.7,33.9 27.2,35.4 22.5,24.3"  fill="#bb0000"/>
    <polygon points="25.8,36.0 22.9,36.7 21.3,24.8"  fill="#ff4400"/>
    <polygon points="21.5,36.9 18.5,36.9 20.0,25.0"  fill="#bb0000"/>
    <polygon points="17.1,36.7 14.2,36.0 18.7,24.8"  fill="#ff4400"/>
    <polygon points="12.8,35.4 10.3,33.9 17.5,24.3"  fill="#bb0000"/>
    <polygon points="9.1,33.0 7.0,30.9 16.5,23.5"    fill="#ff4400"/>
    <polygon points="6.1,29.7 4.6,27.2 15.7,22.5"    fill="#bb0000"/>
    <polygon points="4.0,25.8 3.2,22.9 15.2,21.3"    fill="#ff4400"/>
    <polygon points="3.1,21.5 3.1,18.5 15.0,20.0"    fill="#bb0000"/>
    <polygon points="3.2,17.1 4.0,14.2 15.2,18.7"    fill="#ff4400"/>
    <polygon points="4.6,12.8 6.1,10.3 15.7,17.5"    fill="#bb0000"/>
    <polygon points="7.0,9.1 9.1,7.0 16.5,16.5"      fill="#ff4400"/>
    <polygon points="10.3,6.1 12.8,4.6 17.5,15.7"    fill="#bb0000"/>
    <polygon points="14.2,4.0 17.1,3.2 18.7,15.2"    fill="#ff4400"/>

    <!-- Vertical slit pupil -->
    <ellipse cx="20" cy="20" rx="2.3" ry="9.8" fill="#050000"/>
    <!-- Golden glow rim around slit -->
    <ellipse cx="20" cy="20" rx="3.1" ry="10.5" fill="none"
             stroke="#ffe066" stroke-width="1.1"/>
  </g>
</svg>

<style>
  .sauron-eye {
    flex-shrink: 0;
    overflow: visible;
    cursor: default;
  }
  .sauron-eye.open {
    filter: drop-shadow(0 0 6px #ff550099) drop-shadow(0 0 18px #ff330044);
  }

  /* Scale whole eye body from center */
  .eye-body {
    transform-box: fill-box;
    transform-origin: 50% 50%;
  }

  /* Disconnected: squash flat */
  .sauron-eye:not(.open) .eye-body {
    transform: scaleY(0.05);
    transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* Connected: full open + occasional blink */
  .sauron-eye.open .eye-body {
    animation: blink 8s ease-in-out infinite;
  }
  @keyframes blink {
    0%, 87%          { transform: scaleY(1); }
    90%              { transform: scaleY(0.05); }  /* shut */
    93%, 100%        { transform: scaleY(1); }     /* reopen */
  }

  /* Flicker when active */
  .sauron-eye.open .eye-body {
    filter: none;
  }
</style>
