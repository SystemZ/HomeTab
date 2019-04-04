package pl.systemz.tasktab;

import android.app.Activity;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.widget.Toast;

public class NotificationReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        Integer taskId = intent.getIntExtra("TASK_ID", 0);
        Toast.makeText(context, "Received task " + taskId.toString(), Toast.LENGTH_SHORT).show();

        Intent i = new Intent(context, Task.class);
//        i.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        i.putExtra("TASK_ID", taskId);
        context.startActivity(i);
    }
}
