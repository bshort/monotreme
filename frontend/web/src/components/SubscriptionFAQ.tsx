import Accordion from "@mui/joy/Accordion";
import AccordionDetails from "@mui/joy/AccordionDetails";
import AccordionGroup from "@mui/joy/AccordionGroup";
import AccordionSummary from "@mui/joy/AccordionSummary";

const SubscriptionFAQ = () => {
  return (
    <div className="w-full flex flex-col justify-center items-center">
      <h2 className="text-2xl font-semibold mb-8 dark:text-gray-400">Frequently Asked Questions</h2>
      <AccordionGroup className="w-full max-w-2xl">
        <Accordion>
          <AccordionSummary>Can I use the Free plan in my team?</AccordionSummary>
          <AccordionDetails>
            Of course you can. In the free plan, you can invite up to 5 members to your team. If you need more, you should upgrade to the
            Pro plan.
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary>How many devices can the license key be used on?</AccordionSummary>
          <AccordionDetails>{`It's unlimited for now, but please do not abuse it.`}</AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary>{`Can I get a refund if Monotreme doesn't meet my needs?`}</AccordionSummary>
          <AccordionDetails>
            Yes, absolutely! You can contact us with `yourselfhosted@gmail.com`. I will refund you as soon as possible.
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary>Is there a Lifetime license?</AccordionSummary>
          <AccordionDetails>
            {`As software requires someone to maintain it, so we won't sell a lifetime service, since humans are not immortal yet. But if you
            really want it, please contact us "yourselfhosted@gmail.com".`}
          </AccordionDetails>
        </Accordion>
      </AccordionGroup>
    </div>
  );
};

export default SubscriptionFAQ;
