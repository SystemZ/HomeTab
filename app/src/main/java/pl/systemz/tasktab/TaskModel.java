package pl.systemz.tasktab;

import java.util.List;

public class TaskModel {

    public int id;
    public String name;
    public List<String> tags;
    public int seconds;

    public TaskModel(int id, String name, List<String> tags, int seconds) {
        this.id = id;
        this.name = name;
        this.tags = tags;
        this.seconds = seconds;
    }
}