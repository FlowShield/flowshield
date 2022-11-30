<template>
  <div class="text-center">
    <h1 class="colorful heading text-h1 font-weight-regular mt-10">FlowShield</h1>
<div class="mt-1" >
  <v-container fluid>
    <v-sparkline
      :fill="fill"
      :gradient="selectedGradient"
      :line-width="width"
      :padding="padding"
      :smooth="radius || false"
      :value="value"
      auto-draw
    ></v-sparkline>

    <v-divider></v-divider>
  </v-container>

</div>
<p class="subtitle mt-0 font-weight-regular text-h4">Build a global web3 decentralized private retrieval network for data, building cyber sovereignty</p>

    <div class="mt-15">
      <v-btn x-large color="primary" href="https://www.flowshield.xyz/flowshield_whitepaper" target="_blank">
        <v-icon class="mr-5">mdi-dialpad</v-icon>
        WhitePaper
      </v-btn>
      <v-btn x-large class="ml-15" href="https://github.com/FlowShield/flowshield/tree/main/provider" target="_blank">
        <v-icon class="mr-5">mdi-arrow-up-bold-box-outline</v-icon>
        Deploy miners
      </v-btn>
    </div>

  </div>
</template>
<style lang="scss" scoped>
  .colorful {
    background: linear-gradient(to right, purple, cyan);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
    .subtitle {
    background: linear-gradient(to right, white, cyan);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

</style>

<script>
const gradients = [
  ['#222'],
  ['#42b3f4'],
  ['red', 'orange', 'yellow'],
  ['purple', 'violet'],
  ['#00c6ff', '#F0F', '#FF0'],
  ['#f72047', '#ffd200', '#1feaea']
]

export default {
  data: () => ({
    fill: true,
    selectedGradient: gradients[4],
    gradients,
    padding: 8,
    radius: 10,
    value: [0, 2, 5, 9, 5, 10, 3, 5, 0, 0, 1, 8, 2, 9, 0],
    width: 2
  })
}
</script>
