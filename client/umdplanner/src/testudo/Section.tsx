import Time from './Time';

export class Section{
    code: string;
    times: Array<Time>;
    professor: string;

    conflicts(o: Section):boolean{
        for (let time of this.times){
            for (let otherTime of o.times){
                if (time.conflicts(otherTime)){
                    return true;
                }
            }
        }
        return false;
    }
}

export default Section;