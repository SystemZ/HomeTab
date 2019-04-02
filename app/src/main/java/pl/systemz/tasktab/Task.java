package pl.systemz.tasktab;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import android.util.Log;
import android.widget.EditText;
import android.widget.TextView;

public class Task extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_task);

        Integer taskId = getIntent().getIntExtra("TASK_ID",0);
//        Log.v("tasktab",taskId.toString());

        EditText editText = (EditText)findViewById(R.id.taskIdInput);
        editText.setText(taskId.toString(), TextView.BufferType.EDITABLE);
    }
}
