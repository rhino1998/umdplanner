
export class Duration{
    public start: Date;
    public end: Date;
    conflicts(o: Duration): boolean{
        return (this.start.getTime()<=o.start.getTime() && this.end.getTime()>=o.start.getTime()) ||
        (o.start.getTime()<=this.end.getTime() && o.end.getTime()>=this.start.getTime());
    }
}

export default Duration;