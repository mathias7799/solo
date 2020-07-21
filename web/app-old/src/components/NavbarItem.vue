<template>
  <li class="nav-item" :class="[isSelected ? 'active' : '']">
    <a class="nav-link" href="#">{{ name }}</a>
  </li>
</template>

<script>
module.exports = {
  props: {
    name: String,
    isSelected: { type: Boolean, default: false }
  }
};
</script>