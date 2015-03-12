#include <comedilib.h>

#include "io.h"
#include "channels.h"


static comedi_t     *it_g           = NULL;
static ElevatorType elevatorType    = ET_comedi;

void simulation_elevator_start();
void simulation_dio_write(int channel, int value);
int simulation_dio_read(int channel);
void simulation_data_write(int channel, int value);
int simulation_data_read(int channel);


int io_init(int temp){
	int i = 0;

    elevatorType = temp;
    
    switch(elevatorType){
    case 1: {
        int status = 0;

        it_g = comedi_open("/dev/comedi0");
      
        if (it_g == NULL)
            return 0;

        for (i = 0; i < 8; i++) {
            status |= comedi_dio_config(it_g, PORT1, i,     COMEDI_INPUT);
            status |= comedi_dio_config(it_g, PORT2, i,     COMEDI_OUTPUT);
            status |= comedi_dio_config(it_g, PORT3, i+8,   COMEDI_OUTPUT);
            status |= comedi_dio_config(it_g, PORT4, i+16,  COMEDI_INPUT);
        }

        return (status == 0);
    }
    
    case 0:
        simulation_elevator_start();
        return 1;
        
    default:
        return 0;
    }
}



void io_set_bit(int channel){
    switch(elevatorType){
    case 1:
        comedi_dio_write(it_g, channel >> 8, channel & 0xff, 1);
        break;
        
    case 0:
        simulation_dio_write(channel, 1);
        break;
        
    default:
        break;
    }
}



void io_clear_bit(int channel){
    switch(elevatorType){
    case 1:
        comedi_dio_write(it_g, channel >> 8, channel & 0xff, 0);
        break;
        
    case 0:
        simulation_dio_write(channel, 0);
        break;
        
    default:
        break;
    }
}



void io_write_analog(int channel, int value){
    switch(elevatorType){
    case 1:
        comedi_data_write(it_g, channel >> 8, channel & 0xff, 0, AREF_GROUND, value);
        break;
        
    case 0:
        simulation_data_write(channel, value);
        break;
        
    default:
        break;
    }
}



int io_read_bit(int channel){
    switch(elevatorType){
    case 1: {
        unsigned int data = 0;
        comedi_dio_read(it_g, channel >> 8, channel & 0xff, &data);

        return (int)data;
    }
    case 0:
        return simulation_dio_read(channel);
        
    default:
        return 0;

    }

}



int io_read_analog(int channel){
    switch(elevatorType){
    case 1: {
        lsampl_t data = 0;
        comedi_data_read(it_g, channel >> 8, channel & 0xff, 0, AREF_GROUND, &data);

        return (int)data;
    }
        
    case 0:
        return simulation_data_read(channel);
        
    default:
        return 0;
    }

}




