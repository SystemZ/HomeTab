package pl.systemz.tasktab.model;

import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

public class MqttMsg {

    @SerializedName("type")
    @Expose
    private String type;
    @SerializedName("msg")
    @Expose
    private String msg;
    @SerializedName("id")
    @Expose
    private Integer id;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

}