<template>

  <div id="app">
    <el-tabs v-model="activeName" type="card" @tab-click="handleClick">

      <el-tab-pane name="config">
        <span slot="label"><i class="el-icon-setting"></i> 配置</span>
        <felix-config ref="config"></felix-config>
      </el-tab-pane>

      <el-tab-pane name="machine">
        <span slot="label"><i class="el-icon-s-grid"></i> 机器</span>
        <felix-machine ref="machine"></felix-machine>
      </el-tab-pane>


      <el-tab-pane name="user">
        <span slot="label"><i class="el-icon-user-solid"></i> 用户</span>
        <felix-user ref="user"></felix-user>
      </el-tab-pane>


    </el-tabs>

  </div>

</template>

<script>

import FelixConfig from "./parts/FelixConfig";
import FelixMachine from "./parts/FelixMachine";
import FelixUser from "./parts/FelixUser";

export default {
  name: 'app',
  components: {FelixMachine, FelixConfig,FelixUser},
  data() {
    return {
      activeName: 'machine',
    };
  },
  mounted() {

    for (const one of ['config', 'machine', 'user']) {
      if (window.location.pathname.endsWith(one)){
        this.activeName = one;
        return
      }
    }


  },
  created() {

  },
  methods: {
    handleClick() {
      this.$refs[this.activeName].fetch();
      window.document.title = this.activeName;
    },

  }
};
</script>
<style >
  .search-bar{
    display: flex;
    justify-content: flex-end;
    align-items: center;
  }
  .search-bar > .el-button {
    margin-left: 1rem;
  }
</style>
