<template>
  <a :onclick="onclick">{{ label }}</a>
</template>

<script>
export default {
  name: "NavbarItem",
  props: {
    label: String,
    onclick: String
  }
};
</script>