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
      <el-table-column prop="name" label="name"></el-table-column>
      <el-table-column prop="host" label="host"></el-table-column>
      <el-table-column prop="port" label="port"></el-table-column>
      <el-table-column prop="protocol" label="Protocol"></el-table-column>
      <el-table-column prop="user" label="user"></el-table-column>

      <el-table-column fixed="right" label="Action">
        <template slot-scope="scope">
          <el-button-group type="mini">
            <el-button
                title="webSsh"
                @click="showWebSsh(scope.row)"
                type="primary"
                size="mini"
                icon="el-icon-s-platform"
            ></el-button>
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

        <el-form-item label="host">
          <el-input v-model.trim="form.host" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Port">
          <el-input v-model.number="form.port" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Protocol">
          <el-radio-group v-model="form.protocol">
            <el-radio label="ssh">SSH</el-radio>
            <el-radio label="rdp">RDP(windows)</el-radio>
            <el-radio label="vnc">VNC(macOS/linux)</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="User">
          <el-input v-model.trim="form.user" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="Password">
          <el-input v-model.trim="form.password" autocomplete="off" show-password></el-input>
        </el-form-item>

        <el-form-item label="PrivateKey">
          <el-input type="textarea" v-model.trim="form.private_key" autocomplete="off"></el-input>
        </el-form-item>

        <el-form-item label="PrivateKeyPassword">
          <el-input v-model.trim="form.private_key_password" autocomplete="off" show-password></el-input>
        </el-form-item>

        <el-form-item label="Remark">
          <el-input type="textarea" v-model="form.remark" autocomplete="off"></el-input>
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
  name: "FelixMachine",
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
      axios.get("api/machine", {params: {q}}).then(({status, data}) => {
            if (status === 200) {
              this.tableData = data
            } else {
              this.$message.error(data);
            }
          }
      )
    },
    showWebSsh(row) {
      window.open(row.web_ssh_url, '_blank').focus();
    },
    showDelete(row) {
      this.remove(row)
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
      axios.post("api/machine", this.form).then(({status, data}) => {
        if (status !== 200) {
          this.$message.error(data);
        } else {
          this.dialogFormVisible = false;
          this.fetch()
        }
      })
    },
    remove() {
      axios.delete("api/machine", this.form).then(({status, data}) => {
        if (status !== 200) {
          this.$message.error(data);
        } else {
          this.dialogFormVisible = false;
          this.fetch()
        }
      })
    },
    update() {
      axios.patch("api/machine", this.form).then(({status, data}) => {
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