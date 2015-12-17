require.config({baseUrl:"/assets/js",paths:{"bootstrap-datepicker":"dev/lib/bootstrap-datepicker"}}),define("datetimePick",["util","bootstrap-datepicker"],function(e){$.fn.datepicker.dates["zh-CN"]={days:["星期日","星期一","星期二","星期三","星期四","星期五","星期六","星期日"],daysShort:["周日","周一","周二","周三","周四","周五","周六","周日"],daysMin:["日","一","二","三","四","五","六","日"],months:["一月","二月","三月","四月","五月","六月","七月","八月","九月","十月","十一月","十二月"],monthsShort:["一月","二月","三月","四月","五月","六月","七月","八月","九月","十月","十一月","十二月"],today:"今日",format:"yyyy年mm月dd日",weekStart:1};var t=function(t){this.element=t,this.init=function(t,n,r,i){var s={mode:0,language:"en",format:"yyyy-mm-dd",onChange:null};this.options=$.extend(s,i);var o=$(this.element),u=this;this.picker=o.datepicker({language:this.options.language,autoclose:!0,startDate:n,endDate:r,format:this.options.format}),o.val(e.date2str(t));var a,f,l='<em class="fa fa-calendar"></em>';this.options.mode>0&&(l+='<div style="display:inline;"><select class="form-control cmbhour">');if(this.options.mode>=1){for(a=0;a<24;a++)f=a.toString(),l+="<option"+(a==t.getHours()?' selected="selected"':"")+' value="'+f+'">'+f+"</option>";l+="</select>时"}if(this.options.mode==2){l+='<select class="form-control cmbmin">';for(a=0;a<12;a++)f=(a*5).toString(),l+="<option"+(a*5<=t.getMinutes()&&t.getMinutes()<(a+1)*5?' selected="selected"':"")+' value="'+f+'">'+f+"</option>";l+="</select>分"}l+="</div>";var c=$(l);return this.options.onChange&&(this.picker.on("changeDate",function(){u.options.onChange.call(o)}),c.find("select").change(function(){u.options.onChange.call(o)})),o.after(c),this},this.getDate=function(){var e=$(this.element),t=this.picker.datepicker("getDate"),n=e.next().next();switch(this.options.mode){case 1:t.setHours(parseInt(n.find("select.cmbhour").val()));break;case 2:t.setHours(parseInt(n.find("select.cmbhour").val())),t.setMinutes(parseInt(n.find("select.cmbmin").val()))}return t}};$.fn.gosDatePicker=function(e,n,r,i){return e||(e=new Date),this.each(function(){var s=$(this),o=s.data("gos.datepicker");if(o)return;s.data("gos.datepicker",(new t(this)).init(e,n,r,i))})}});