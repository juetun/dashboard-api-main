(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-b6eb1da8"],{9561:function(e,t,a){"use strict";a.r(t);var r=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",[a("Card",[a("div",[a("Form",{ref:"formValidate",attrs:{"label-position":"left",model:e.formValidate,rules:e.ruleValidate,"label-width":120}},[a("FormItem",{attrs:{label:"Name",prop:"name"}},[a("Input",{attrs:{placeholder:"title"},model:{value:e.formValidate.name,callback:function(t){e.$set(e.formValidate,"name",t)},expression:"formValidate.name"}})],1),a("FormItem",{attrs:{label:"DisplayName",prop:"displayName"}},[a("Input",{attrs:{placeholder:"title"},model:{value:e.formValidate.displayName,callback:function(t){e.$set(e.formValidate,"displayName",t)},expression:"formValidate.displayName"}})],1),a("FormItem",{attrs:{label:"SeoDescription",prop:"seoDescription"}},[a("Input",{attrs:{type:"textarea",autosize:{minRows:2},placeholder:"Enter seo description..."},model:{value:e.formValidate.seoDescription,callback:function(t){e.$set(e.formValidate,"seoDescription",t)},expression:"formValidate.seoDescription"}})],1),a("FormItem",[a("Button",{attrs:{type:"primary"},on:{click:function(t){return e.handleSubmit("formValidate")}}},[e._v("Submit")]),a("Button",{staticStyle:{"margin-left":"8px"},on:{click:function(t){return e.handleReset("formValidate")}}},[e._v("Reset")])],1)],1)],1)])],1)},n=[],i=(a("5a39"),a("d28d")),o={data:function(){return{formValidate:{name:"",displayName:"",seoDescription:""},ruleValidate:{name:[{required:!0,message:"The name cannot be empty",trigger:"blur"},{max:100,message:"The name length is too long",trigger:"blur"}],displayName:[{required:!0,message:"The displayName cannot be empty",trigger:"blur"},{max:100,message:"The displayName length is too long",trigger:"blur"}],seoDescription:[{required:!0,message:"The seo description can not be empty",trigger:"blur"},{max:250,message:"The seo description length is too long",trigger:"blur"}]},tagId:0}},mounted:function(){var e=this.$route.query.id;this.tagId=e,this.defaultData(e)},methods:{defaultData:function(e){var t=this;Object(i["d"])(e).then(function(e){t.formValidate.name=e.data.data.Name,t.formValidate.displayName=e.data.data.DisplayName,t.formValidate.seoDescription=e.data.data.SeoDesc}).catch(function(e){t.$Message.error("操作失败"+e)})},handleSubmit:function(e){var t=this,a=this;this.$refs[e].validate(function(e){e?Object(i["e"])(a.tagId,a.formValidate.name,a.formValidate.displayName,a.formValidate.seoDescription).then(function(e){0===e.data.code?(t.$Message.success(e.data.message),setTimeout(function(){t.$router.push("/backend/tag/list")},2e3)):t.$Message.error(e.data.message)}).catch(function(e){t.$Message.error("操作失败"+e)}):t.$Message.error("Fail!")})},handleReset:function(e){this.$refs[e].resetFields()}}},s=o,l=a("620d"),d=Object(l["a"])(s,r,n,!1,null,null,null);t["default"]=d.exports},d28d:function(e,t,a){"use strict";a.d(t,"a",function(){return o}),a.d(t,"d",function(){return s}),a.d(t,"e",function(){return l}),a.d(t,"b",function(){return d}),a.d(t,"c",function(){return m});var r=a("0d5e"),n=Object(r["a"])(),i={headers:{"Content-Type":"multipart/form-data"}};function o(e){return n.get("/console/tag/",{params:e})}function s(e,t){return n.get("/console/tag/edit/"+e,{params:t})}function l(e,t,a,r){return n.put("/console/tag/"+e,{name:t,displayName:a,seoDesc:r},i)}function d(e,t,a){return n.post("/console/tag/",{name:e,displayName:t,seoDesc:a},i)}function m(e,t){return n.delete("/console/tag/"+e,{params:t})}}}]);