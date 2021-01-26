package pl.systemz.hometab.receiver;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.widget.Toast;

import pl.systemz.hometab.Task;

public class NotificationReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        Integer taskId = intent.getIntExtra("TASK_ID", 0);
        String taskName = intent.getStringExtra("TASK_NAME");
        String action = intent.getStringExtra("ACTION");

        if (action.equals("VIEW")) {
            Intent i = new Intent(context, Task.class);
            //i.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
            i.putExtra("TASK_ID", taskId);
            i.putExtra("TASK_NAME", taskName);
            context.startActivity(i);
        } else if (action.equals("STOP")) {
            Toast.makeText(context, "Can't stop, won't stop", Toast.LENGTH_SHORT).show();
        }

    }
}
