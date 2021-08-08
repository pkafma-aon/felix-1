<template>
  <div>
    <el-row align="middle" style="margin-bottom:1rem" justify="start" type="flex">
      <el-col :span="4">

      </el-col>
      <el-col :span="20" class="search-bar">
        <el-input v-model.trim="q"
                  clearable
                  placeholder="模糊搜索用户名"
                  prefix-icon="el-icon-search"
                  size="mini"
                  @change="fetch">
        </el-input>

        <el-button icon="el-icon-search" size="mini"
                   title="搜索/刷新按钮"
                   type="primary"
                   @click="fetch">
        </el-button>
      </el-col>
    </el-row>

    <el-table :data="tableData" border style="width: 100%" stripe>

      <el-table-column prop="name" label="name">
      </el-table-column>
      <el-table-column prop="value" label="value">
      </el-table-column>
      <el-table-column prop="def" label="Default">
      </el-table-column>
      <el-table-column prop="desc" label="Desc">
      </el-table-column>


      <el-table-column fixed="right" label="Action" width="240">
        <template slot-scope="scope">
          <el-link
              title="edit config"
              @click="showForm(scope.row)"
              type="success"
              size="mini"
              icon="el-icon-edit"
          ></el-link>
        </template>
      </el-table-column>
    </el-table>


    <!--edit-->
    <el-dialog title="ssh" :visible.sync="dialogFormVisible">
      <el-form :model="form" label-position="right" label-width="120px" size="mini">
        <el-form-item label="Name">
          <el-input v-model="form.name" autocomplete="off" disabled></el-input>
        </el-form-item>

        <el-form-item label="Default">
          <el-input v-model="form.def" autocomplete="off" disabled></el-input>
        </el-form-item>

        <el-form-item label="Desc">
          <el-input v-model="form.desc" autocomplete="off" disabled></el-input>
        </el-form-item>

        <el-form-item label="Value">
          <el-input v-model="form.value" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" type="danger" size="mini">Close</el-button>
        <el-button type="primary" @click="updateConfig" size="mini">Submit</el-button>
      </div>
    </el-dialog>

  </div>

</template>

<script>
import axios from 'axios';

export default {
  name: "FelixConfig",
  data() {
    return {
      formLabelWidth: "120",
      dialogFormVisible: false,
      total: 0,
      page: 1,
      size: 10,
      tableData: [],
      form: {},
      q: "",
    };
  },
  mounted() {
    this.fetch();
  },
  created() {
  },
  methods: {
    fetch() {
      let q = ''
      axios.get("api/config", {params: {q}}).then(({status, data}) => {
            if (status === 200) {
              this.tableData = data
            } else {
              this.$message.error(data);
            }
          }
      )
    },
    showForm(row) {
      this.form = row;
      this.dialogFormVisible = true
    },
    updateConfig() {
      axios.patch("api/config", this.form).then(({status, data}) => {
        if (status !== 200) {
          this.$message.error(data);
        } else {
          this.dialogFormVisible = false;
          this.fetch()
        }
      })
    },
  }
}
</script>

<style scoped>

</style>