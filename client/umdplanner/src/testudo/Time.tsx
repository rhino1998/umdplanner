import Duration from './Duration';


export class Time{
    room: string;
    duration: Duration;

    conflicts(o:Time): boolean {
        return this.duration.conflicts(o.duration)
    }
}

export default Time;