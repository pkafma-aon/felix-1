<template>
  <div>
    <el-row align="middle" style="margin-bottom:1rem" justify="start" type="flex">
      <el-col :span="4">
        <el-button title="" size="mini" type="success" @click="showCreate">Create</el-button>

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

      <el-table-column prop="id" label="ID"></el-table-column>
      <el-table-column prop="name" label="Name"></el-table-column>
      <el-table-column prop="role" label="Role"></el-table-column>
      <el-table-column prop="account" label="Account"></el-table-column>
      <el-table-column prop="email" label="Email"></el-table-column>
      <el-table-column prop="phone" label="Phone"></el-table-column>
      <el-table-column prop="password" label="Password"></el-table-column>
      <el-table-column prop="github_account" label="GithubAccount"></el-table-column>


      <el-table-column fixed="right" label="Action" width="240">
        <template slot-scope="scope">
          <el-button-group type="mini">
            <el-button
                title="edit config"
                @click="showUpdate(scope.row)"
                type="success"
                size="mini"
                icon="el-icon-edit"
            ></el-button>
            <el-button
                title="remove"
                @click="showDelete(scope.row)"
                type="danger"
                size="mini"
                icon="el-icon-delete"
            ></el-button>
            <el-button
                title="Sync Github Public key"
                @click="doSyncGithubPublicKey(scope.row)"
                type="info"
                size="mini"
                icon="el-icon-refresh"
            ></el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>


    <!--edit-->
    <el-dialog :title="createOrUpdate" :visible.sync="dialogFormVisible">
      <el-form :model="form" label-position="right" label-width="120px" size="mini">
        <el-form-item label="Name">
          <el-input v-model="form.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="Role">
          <el-radio-group v-model="form.role">
            <el-radio label="Reporter">Reporter</el-radio>
            <el-radio label="Developer">Developer</el-radio>
            <el-radio label="Maintainer">Maintainer</el-radio>
            <el-radio label="Admin">Admin</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="Email">
          <el-input v-model.trim="form.email" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Account">
          <el-input v-model="form.account" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Phone">
          <el-input v-model.trim="form.phone" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Github Account">
          <el-input v-model.trim="form.github_account" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="Public Key">
          <el-input type="textarea"
                    :rows="10"
                    v-model.trim="form.public_key" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Password">
          <el-input v-model.trim="form.password" autocomplete="off" show-password></el-input>
        </el-form-item>

      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" type="danger" size="mini">Close</el-button>
        <el-button type="primary" @click="doSubmit" size="mini">Submit</el-button>
      </div>
    </el-dialog>

  </div>

</template>

<script>
import axios from 'axios';

export default {
  name: "FelixUser",
  data() {
    return {
      createOrUpdate: 'create',
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
      axios.get("api/user", {params: {q}}).then(({status, data}) => {
            if (status === 200) {
              this.tableData = data
            }else{
              this.$message.error(data)
            }
          }
      )
    },
    showDelete(row) {
      this.remove(row)
    },
    doSyncGithubPublicKey(row) {
      let uri = "api/user-github-public-key-sync";
      axios.post(uri, row).then(({status, data}) => {
            if (status === 200) {
              this.fetch()
            }else{
              this.$message.error(data)
            }
          }
      )
    },
    showUpdate(row) {
      this.createOrUpdate = 'update'
      this.form = row;
      this.dialogFormVisible = true
    },
    doSubmit() {
      if (this.createOrUpdate === "create") {
        this.create()
      } else {
        this.update()
      }
    },
    showCreate() {
      this.createOrUpdate = 'create'
      this.form = {};
      this.dialogFormVisible = true
    },
    create() {
      axios.post("api/user", this.form).then(({status, data}) => {
        if (status !== 200) {
          this.$message.error(data);
        } else {
          this.dialogFormVisible = false;
          this.fetch()
        }
      })
    },
    remove() {
      axios.delete("api/user", this.form).then(({status, data}) => {
        if (status !== 200) {
          this.$message.error(data);
        } else {
          this.dialogFormVisible = false;
          this.fetch()
        }
      })
    },
    update() {
      axios.patch("api/user", this.form).then(({status, data}) => {
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